import { GetRelativeMoistureRequest } from '@irrigation/protobuf/irrigation_pb';
import { Typography } from '@material-ui/core';
import { Error as GrpcError } from 'grpc-web';
import React, { useEffect, useState } from 'react';
import irrService from '../irrService';


export default function MoistureSensor({ channel }: { channel: number }) {
  const [error, setError] = useState<GrpcError | null>(null);
  const [isLoaded, setIsLoaded] = useState(false);
  const [moisture, setMoisture] = useState(0);

  const getMoisture = () => {
    const req = new GetRelativeMoistureRequest();
    req.setChannel(channel);

    irrService.getRelativeMoisture(req, {}, (err, resp) => {
      if (err) {
        setError(err)
        return;
      }
      setIsLoaded(true);
      setMoisture(resp.getMoisture());
    });
  }

  useEffect(() => {
    getMoisture();
    const interval = setInterval(() => getMoisture(), 5000);
    return () => clearInterval(interval);
  }, []);

  if (error) {
    return <div>Error: {error.message}</div>;
  } else if (!isLoaded) {
    return <div>Loading ...</div>;
  } else {
    return (
      <div>
        <Typography variant="overline" display="block" gutterBottom color="textSecondary">Channel {channel}</Typography>
        <Typography variant="h4" component="p" color="textPrimary">
          {Intl.NumberFormat('en-US', { maximumSignificantDigits: 3 }).format(moisture * 100)}%
        </Typography>
      </div>
    )
  }
}
