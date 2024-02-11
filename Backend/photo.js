// para las fotos photo.js
const mongoose = require('mongoose');

const photoSchema = new mongoose.Schema({
  image: {
    type: String,
    required: true
  },
  uploadDate: {
    type: Date,
    default: Date.now
  }
});

module.exports = mongoose.model('Photo', photoSchema);
