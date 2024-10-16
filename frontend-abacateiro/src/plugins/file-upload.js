import axios from 'axios';

import { LocalStorage } from "quasar";

const BASE_URL = 'http://localhost:8888';

function uploadXML(formData, api) {
  const token = LocalStorage.getItem("token");
  const config = {
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "multipart/form-data",
    },
  };

  const url = `${BASE_URL}` + api;
  return axios
    .post(url, formData, config)
    .then((x) => x.data)
    .catch((err) => {
      throw err;
    });
}

function uploadCSV(formData, api) {
  const token = LocalStorage.getItem("token");
  const config = {
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "multipart/form-data",
    },
  };

  const url = `${BASE_URL}` + api;
  return axios
    .post(url, formData, config)
    .then((x) => x.data)
    .catch((err) => {
      throw err;
    });
}

function uploadData(formData, api) {
  const token = LocalStorage.getItem("token");
  const config = {
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "multipart/form-data",
    },
  };

  const url = `${BASE_URL}` + api;
  return axios
    .post(url, formData, config)
    .then((x) => x.data)
    .catch((err) => {
      throw err;
    });
}

export { uploadCSV, uploadXML, uploadData };
