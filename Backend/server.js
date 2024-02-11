const express = require('express');
const multer = require('multer');
const path = require('path');
const mongoose = require('mongoose');
const Photo = require('./photo'); // Asegúrate de que la ruta sea correcta
const cors = require('cors');

const app = express();
const port = 3001;

// Conexión a MongoDB (ajusta tu cadena de conexión)
mongoose.connect('mongodb://localhost:27017/photo', {
  useNewUrlParser: true,
  useUnifiedTopology: true
})
.then(() => console.log('Connected to MongoDB'))
.catch(err => console.error('MongoDB connection error:', err));

// Configuración de almacenamiento de Multer
const storage = multer.diskStorage({
  destination: './uploads/',
  filename: function (req, file, cb) {
    cb(null, file.fieldname + '-' + Date.now() + path.extname(file.originalname));
  }
});

// Configuración de filtro de archivos de Multer
const upload = multer({ 
  storage: storage,
  fileFilter: (req, file, cb) => {
    if (!file.originalname.match(/\.(jpg|jpeg|png|gif)$/)) {
      return cb(new Error('Only image files are allowed!'), false);
    }
    cb(null, true);
  }
});

// Middleware de CORS
app.use(cors());

// Ruta para la carga de fotos
app.post('/uploadPhotos', upload.array('photos'), async (req, res) => {
  try {
    const files = req.files;

    if (!files || files.length === 0) {
      return res.status(400).json({ message: 'No files uploaded' });
    }

    const uploadedPhotos = [];
    for (const file of files) {
      const newPhoto = new Photo({
        image: file.path, // Asegúrate de que 'image' coincida con tu esquema de Mongoose
        uploadDate: new Date()
      });
      await newPhoto.save();
      uploadedPhotos.push(newPhoto);
    }

    res.status(201).json({ message: 'Photos uploaded successfully', photos: uploadedPhotos });
  } catch (error) {
    console.error('Error uploading photos:', error);
    res.status(500).json({ error: error.message || 'Failed to upload photos' });
  }
});

// Nueva ruta para obtener las fotos de la base de datos
app.get('/getPhotos', async (req, res) => {
  try {
    const photos = await Photo.find(); // Obtener todas las fotos de la base de datos
    res.json(photos);
  } catch (error) {
    console.error('Error retrieving photos:', error);
    res.status(500).json({ error: 'Failed to retrieve photos' });
  }
});

// Iniciar el servidor
app.listen(port, () => {
  console.log(`Server running at http://localhost:${port}`);
});
