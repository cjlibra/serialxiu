var SerialPort = require("serialport").SerialPort;

var serialPort = new SerialPort('COM3',{baudrate: 38400}, true);

var scanstr="\x02\x01\x08\x00\x00\x01\x00\x0a";

var program = require('commander');

program
  .version('0.0.1')
  .option('-r, --read', 'read  tag data')
  .option('-w, --write', 'write  tag data')
  .option('-i, --id <id>', 'term id')
  .option('-f, --file <filename>', 'Add file name for read write')
  .parse(process.argv);

console.log('you ordered a pizza with:');
if (program.read && program.id && program.file) {
	console.log('  read tag data from tag id');
}
if (program.write && program.id && program.file) {
	console.log('  write tag data from tag id');
}
 
 
 
 
 serialPort.on ('open', function (error) {
     if ( error ) {
        console.log('failed to open: '+error);
      } else {
        console.log('open');
	
	
		serialPort.write(scanstr, function(err, results) {
			if (error){
				console.log("failed write")
			    console.log('err ' + err);
		        console.log('results ' + results);
			}else{
		        
			}
		});
		 
	serialPort.on ('data', function( data ) {
		console.log("data" + data.toString());
			 
	});	
		
	   
  }
});