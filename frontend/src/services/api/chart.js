'use strict'

import axios from 'axios';

const API_URL = 'http://localhost:8080/api/v1/chart';
const headers =  { 
  'Content-Type': 'application/json'
}

const getChart = () => {
  const config = {
    method: 'get',
    url: API_URL,
    headers,
  };

  return axios.request(config)
}

const setDeviceInChart = (device, data) => {
  const config = {
    method: 'post',
    url: API_URL + '/' + device,
    headers,
    data
  };

  return axios.request(config)
}

export default {
  getChart,
  setDeviceInChart,
}
