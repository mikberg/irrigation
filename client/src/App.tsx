import { Container, Grid, Paper, Typography } from '@material-ui/core';
import { createStyles, makeStyles, Theme } from '@material-ui/core/styles';
import React from 'react';
import './App.css';
import MoistureSensor from './components/MoistureSensor';
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
              <Paper className={classes.paper}>xs</Paper>
            </Grid>
          </Grid>


          <Typography variant="h6" gutterBottom>Apply water</Typography>

          <Grid container spacing={1}>
            <Grid item xs>
              <WaterButton channel={0} />
            </Grid>
            <Grid item xs>
              <WaterButton channel={1} />
            </Grid>
            <Grid item xs>
              <WaterButton channel={2} />
            </Grid>
          </Grid>


          <Paper className={classes.paper}>
            <WaterLevelSensor />
          </Paper>

        </Container>
      </main>
    </div>
  )
}

export default App
