import * as irr_pb from '@irrigation/protobuf/irrigation_pb';
import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import './index.css';
import irrService from './irrService';

const waterRequest = new irr_pb.WaterRequest();
waterRequest.setDuration(10);
waterRequest.setChannel(0);

irrService.water(waterRequest, {}, (err, response) => {
  if (err) {
    console.log(err);
    return;
  }
  console.log(response);
});

ReactDOM.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
  document.getElementById('root')
)
