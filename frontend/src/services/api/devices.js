'use strict'

import axios from 'axios';

const API_URL = 'http://localhost:8080/api/v1/devices';
const headers =  { 
  'Content-Type': 'application/json'
}

const getDevices = () => {
  const config = {
    method: 'get',
    url: API_URL,
    headers,
  };

  return axios.request(config)
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

const getDeviceById = (deviceId) => {
  const config = {
    method: 'get',
    url: API_URL + '/' + deviceId,
    headers,
  };

  return axios.request(config)
}

const getRoute = (sourceId, targetId, type = "distance") => {
  const config = {
    method: 'get',
    url: API_URL + '/route/' + sourceId + '/' + targetId + '?type=' + type,
    headers,
  };

  return axios.request(config)
}

const deleteDevice = (deviceLabel) => {
  const config = {
    method: 'delete',
    url: API_URL + '/' + deviceLabel,
    headers,
  };

  return axios.request(config)
}

const sendMessage = (data) => {
  const config = {
    method: 'post',
    url: API_URL + '/requests',
    headers,
    data,
  };

  return axios.request(config)
}

export default {
  getDevices,
  insertDevice,
  getDeviceById,
  getRoute,
  deleteDevice,
  sendMessage,
}
