// Socket IO Client
var io = require('socket.io-client');
var socket = io.connect('http://localhost:5000');

socket.on('data', function(data){
  console.log(data);
});
