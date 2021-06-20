import { GetWaterLevelRequest, GetWaterLevelResponse } from '@irrigation/protobuf/irrigation_pb';
import { Typography } from '@material-ui/core';
import { Error as GrpcError } from 'grpc-web';
import React, { useEffect, useState } from 'react';
import irrService from '../irrService';

export default function WaterLevelSensor() {
  const [error, setError] = useState<GrpcError | null>(null);
  const [isLoaded, setIsLoaded] = useState(false);
  const [response, setResponse] = useState<GetWaterLevelResponse | null>(null);

  const getWaterLevel = () => {
    const req = new GetWaterLevelRequest();

    irrService.getWaterLevel(req, {}, (err, resp) => {
      if (err) {
        setError(err);
        return;
      }
      setError(null);
      setResponse(resp);
    });
  };

  useEffect(() => {
    getWaterLevel();

    const interval = setInterval(() => getWaterLevel(), 5000);
    return () => clearInterval(interval);
  }, []);

  if (error) {
    return <div>Error: {error.message}</div>;
  } else if (!response) {
    return <div>Loading ...</div>;
  } else {
    return (
      <div>
        <Typography variant="overline" display="block" gutterBottom color="textSecondary">Water</Typography>
        <Typography variant="h4" component="p" color="textPrimary">
          {Intl.NumberFormat('en-US', { maximumSignificantDigits: 4 }).format(response.getLiters())}L
        </Typography>
      </div>
    );
  }
}
