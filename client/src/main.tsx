import blupp from '@irrigation/protobuf/blupp';
import * as car_pb from '@irrigation/protobuf/car_pb';
import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import './index.css';


blupp('world');

const car = new car_pb.Car();
car.setMake('Peugeot');
car.setModel('307 SWs');
console.log(car.toObject());

ReactDOM.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
  document.getElementById('root')
)
