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

const getRoute = (sourceId, targetId) => {
  const config = {
    method: 'get',
    url: API_URL + '/route/' + sourceId + '/' + targetId,
    headers,
  };

  return axios.request(config)
}

export default {
  insertDevice,
  getRoute,
}
