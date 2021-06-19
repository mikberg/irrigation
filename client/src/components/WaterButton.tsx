import { WaterRequest } from '@irrigation/protobuf/irrigation_pb';
import { Button, Snackbar } from '@material-ui/core';
import { Error as GrpcError } from 'grpc-web';
import React, { useState } from 'react';
import irrService from '../irrService';

export default function WaterButton({ channel }: { channel: number }) {
  const [error, setError] = useState<GrpcError | null>(null);
  const [isWatering, setIsWatering] = useState(false);
  const [isDone, setIsDone] = useState(false);


  const waterTimer = React.useRef<NodeJS.Timer>();
  const successTimer = React.useRef<NodeJS.Timer>();

  const handleButtonClick = () => {
    if (isWatering) {
      return;
    }

    setIsWatering(true);

    const req = new WaterRequest();
    req.setChannel(channel);
    req.setDuration(2);

    irrService.water(req, {}, (err, resp) => {
      if (err) {
        setError(err);
        setIsWatering(false);
        setIsDone(true);
      } else {
        setError(null);
      }

      waterTimer.current = setTimeout(() => {
        setIsWatering(false);
        setIsDone(true);
      }, req.getDuration() * 1000);
    });
  };

  const handleSuccessClose = () => {
    setIsDone(false);
  };

  return (<div>
    <Button key="button" variant="contained" color="primary" disabled={isWatering} onClick={handleButtonClick}>
      Water Channel {channel}
    </Button>
    <Snackbar key="snackbar" open={isDone} autoHideDuration={2000} onClose={handleSuccessClose} message={error ? `Error! ${error.message}` : 'Watering complete!'} />
  </div>)
}
