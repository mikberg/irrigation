import blupp from '@irrigation/protobuf/blupp';
import * as irr_pb from '@irrigation/protobuf/irrigation_pb';
import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import './index.css';


blupp('world');

const waterRequest = new irr_pb.WaterRequest();
waterRequest.setDuration(10);
waterRequest.setChannel(0);
console.log(waterRequest.toObject());

ReactDOM.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
  document.getElementById('root')
)
