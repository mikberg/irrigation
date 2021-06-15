#!/usr/bin/env node

const { createServer } = require('vite');
const { handleHMRUpdate } = require('vite/dist/node');
const path = require('path');
const { EventEmitter } = require('stream');
const { readFileSync, readFile } = require('fs');

console.log('hello from prodev');

process.chdir('client/vite')
console.log('cwd', process.cwd());

class MyWatcher extends EventEmitter {
  constructor() {
    super();

    this.files = [];
    this.server = null;
  }

  setServer(server) {
    this.server = server;
  }

  add(paths, origAdd, internal) {
    console.log('[MyWatcher] Added file!', paths)

    this.files.push(paths);

    return this;
  }

  async restartServer() {
    const newServer = await restartServer();

    await this.server.close();

    this.server = newServer;
    await this.server.listen();
    console.log('server restarted...');
  }

  async blupp() {
    console.log('Emitting change!!!', this.files);
    this.emit('change', 'somepath.ts');

    const { ws, config, moduleGraph } = this.server;

    // Reloads the page, but doesn't change content
    // ws.send({
    //   type: 'full-reload',
    //   path: '*'
    // });

    // Works!!
    this.restartServer();

    // console.log('contents', readFileSync('client/src/App.tsx').toString());
    // const mods = moduleGraph.getModulesByFile('src/App.tsx');
    // console.log('mods', mods);

    //   const hmrContext = {
    //     file: 'src/App.tsx',
    //     timestamp: Date.now(),
    //     modules: mods ? [...mods] : [],
    //     read: () => readFileSync('client/src/App.tsx'),
    //     server: this.server,
    //   };
    //   for (const plugin of config.plugins) {
    //     if (plugin.handleHotUpdate) {
    //       const filteredModules = await plugin.handleHotUpdate(hmrContext);
    //       if (filteredModules) {
    //         hmrContext.modules = filteredModules;
    //       }
    //     }
    //   }
    //   console.log(hmrContext);
  }
}

const myWatcher = new MyWatcher();

process.stdin.on('data', data => {
  if (data.toString().trim() == 'IBAZEL_BUILD_COMPLETED SUCCESS') {
    console.log('Files were updated! :---) ')
    myWatcher.blupp();
  } else {
    console.log('something else happened');
  }
});

const restartServer = async () => {
  return await createServer({
    configFile: '../vite.config.ts',
    root: path.join(__dirname, '..'),
    server: {
      port: 1337,
    },
  });
};

(async () => {
  console.log('dirname', __dirname);

  const server = await restartServer();

  server.watcher = myWatcher;
  myWatcher.setServer(server);

  console.log(server.watcher);
  await server.listen();
})()
