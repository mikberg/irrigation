import { Slider } from '@material-ui/core';
import React from 'react';

const marks = [
  { label: '5s', value: 5 },
  { label: '10s', value: 10 },
  { label: '20s', value: 20 },
  { label: '30s', value: 30 },
];

export default function WaterSlider({ onChange }: { onChange: (value: number) => void }) {
  const handleChange = (event: any, newValue: number | number[]) => {
    if (!!onChange) {
      if (Array.isArray(newValue)) {
        onChange(newValue[0]);
      } else {
        onChange(newValue);
      }
    }
  }

  return (
    <Slider
      defaultValue={10}
      step={5}
      marks={marks}
      min={5}
      max={30}
      valueLabelDisplay="auto"
      onChange={handleChange}
    />
  )
}
