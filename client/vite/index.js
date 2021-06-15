#!/usr/bin/env node

const { createServer } = require('vite');
const path = require('path');
const { EventEmitter } = require('stream');

class BazelWatcher extends EventEmitter {
  constructor(stream) {
    super();

    stream.on('data', (data) => this.onData(data))
  }

  onData(data) {
    const event = data.toString().trim();

    switch (event) {
      case 'IBAZEL_BUILD_STARTED':
        this.emit('started');
        break;
      case 'IBAZEL_BUILD_COMPLETED SUCCESS':
        this.emit('build_success');
        break;
      case 'IBAZEL_BUILD_COMPLETED FAILURE':
        this.emit('build_failure');
        break;
      default:
        throw new Error(`Unknown event from ibazel: ${event}`)
    }
  }
}


class Coordinator {
  constructor(config) {
    this.config = config;

    this.bazelWatcher = new BazelWatcher(process.stdin);
    this.bazelWatcher.on('build_success', this.onBuildSuccess.bind(this));
  }

  async createServer() {
    const server = await createServer(this.config);

    // server.watcher.on('change', (filePath) => {
    //   console.log('FILE CHANGED; RESTARTING SERVER!!', path.relative(this.config.root, filePath));
    //   this.restartServer();
    // });

    // server.watcher.add([
    //   '../protobuf/car_pb.js',
    //   '../protobuf/car_pb.mjs',
    //   '../protobuf/car_pb.d.ts',
    // ]);

    // server.watcher
    //   .on('add', path => console.log(`File ${path} has been added`))
    //   .on('change', path => console.log(`File ${path} has been changed`))
    //   .on('unlink', path => console.log(`File ${path} has been removed`));

    server.watcher.on('change', filePath => {
      const relativePath = path.relative(this.config.root, filePath);
      console.log(`File ${relativePath} changed`);
      if (relativePath.startsWith('..')) {
        console.info('')
        this.restartServer();
      }
    });

    // // More possible events.
    // server.watcher
    //   .on('addDir', path => console.log(`Directory ${path} has been added`))
    //   .on('unlinkDir', path => console.log(`Directory ${path} has been removed`))
    //   .on('error', error => console.log(`Watcher error: ${error}`))
    //   .on('ready', () => console.log('Initial scan complete. Ready for changes'))
    //   .on('raw', (event, path, details) => { // internal
    //     console.log('Raw event info:', event, path, details);
    //   });

    // console.log('WATCHERS', server.watcher.getWatched());

    return server;
  }

  async start() {
    this.server = await this.createServer();
    await this.server.listen();
  }

  async onBuildSuccess() {
    // console.log('WATCHERS', this.server.watcher.getWatched());
    // await this.fullRefresh();
    // await this.restartServer();
  }

  // Restarts the server, causing the client to reload
  async restartServer() {
    const newServer = await this.createServer(this.config);

    await this.server.close();

    this.server = newServer;
    await this.server.listen();
  }

  // Causes client to refresh
  async fullRefresh() {
    const { ws, config, moduleGraph } = this.server;

    ws.send({
      type: 'full-reload',
      path: '*'
    });
  }
}

const coordinator = new Coordinator({
  configFile: path.join(process.env.VITE_ROOT, 'vite.config.ts'),
  root: process.env.VITE_ROOT,
  server: {
    port: 3000,
    watch: {
      usePolling: true,
    }
  }
});

(async () => await coordinator.start())();
