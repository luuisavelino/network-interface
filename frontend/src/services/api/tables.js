'use strict'

import axios from 'axios';

const API_URL = 'http://localhost:8080/api/v1/tables';
const headers =  { 
  'Content-Type': 'application/json'
}

const getTable = () => {
  const config = {
    method: 'get',
    url: API_URL,
    headers,
  };

  return axios.request(config)
}

export default {
  getTable,
}
