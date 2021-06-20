import { Box, Card, CardContent, Container, Grid, Paper } from '@material-ui/core';
import { createStyles, makeStyles, Theme } from '@material-ui/core/styles';
import React, { useState } from 'react';
import './App.css';
import MoistureSensor from './components/MoistureSensor';
import TimeSlider from './components/TimeSlider';
import WaterButton from './components/WaterButton';
import WaterLevelSensor from './components/WaterLevelSensor';

const useStyles = makeStyles((theme: Theme) => createStyles(
  {
    root: {
      flexGrow: 1,
    },
    content: {
      flexGrow: 1,
      overflow: 'auto',
    },
    paper: {
      padding: theme.spacing(1),
      textAlign: 'center',
    }
  }
));

function App() {
  const [waterDuration, setWaterDuration] = useState<number>(10);

  const classes = useStyles();
  return (
    <div className={classes.root}>
      <main className={classes.content}>
        <Container maxWidth="lg" style={{ padding: 20 }}>
          <Grid container spacing={1}>
            <Grid item xs>
              <Paper className={classes.paper}>
                <MoistureSensor channel={0} />
              </Paper>
            </Grid>
            <Grid item xs>
              <Paper className={classes.paper}>
                <MoistureSensor channel={1} />
              </Paper>
            </Grid>
            <Grid item xs>
              <Paper className={classes.paper}>
                <MoistureSensor channel={2} />
              </Paper>
            </Grid>
          </Grid>



          <Box my={2} p={0}>
            <Card>
              <CardContent>
                <Box my={1}>
                  <WaterButton channel={0} duration={waterDuration} />
                </Box>
                <Box my={1}>
                  <WaterButton channel={1} duration={waterDuration} />
                </Box>
                <Box my={1}>
                  <WaterButton channel={2} duration={waterDuration} />
                </Box>

                <TimeSlider onChange={setWaterDuration} />
              </CardContent>
            </Card>
          </Box>

          <Paper className={classes.paper}>
            <WaterLevelSensor />
          </Paper>

        </Container>
      </main>
    </div >
  )
}

export default App
