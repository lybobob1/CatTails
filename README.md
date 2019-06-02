# CatTails

## Overview  
This project is a redteam implant that leverages raw sockets to  
send/recieve callbacks from a C2 infrastructure.  
  
The callbacks and communication happen over UDP port 53. You will be able  
to send/execute commands on the remote host via a C2 server.  

### Features

- ##### custom commands*
    - Ex. If you want to flush `iptables` CatTails will provide an
      abstraction for you to do this. Instead of sending the shell commands
      necessary you will be able to run a CatTails command like `drop-rules`. 
- ##### Command feedback/output*
    - CatTails will send you the output of a command (if there is any)  
      and let you know if your command completed successfully.  

(*) Work in progress  
(x) Completed

- #### More coming soon!
