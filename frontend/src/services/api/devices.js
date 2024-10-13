'use strict'

import axios from 'axios';

const API_URL = 'http://localhost:8080/api/v1/devices';
const headers =  { 
  'Content-Type': 'application/json'
}

const insertDevice = (data) => {
  const config = {
    method: 'post',
    url: API_URL,
    headers,
    data,
  };

  return axios.request(config)
}

export default {
  insertDevice,
}
