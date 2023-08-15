# task-bomberman-dom
grit:lab Åland Islands 2023

After run project, open the browser on `http://localhost:3000/`  

## docker run  

To show two docker containers inside docker desktop or terminal, one for backend, and one for frontend (docker desktop must be started before), inside root folder of repository where `run.sh` placed, using terminal execute command:
- `./run.sh`  
to show containers in terminal:  
- `docker ps -a`  

## dev run  

### backend  

- `cd backend`
- `go run .`

### frontend  

- `cd frontend`
- `npm i`
- `npm run start`

## task description and audit questions, on github  

https://github.com/01-edu/public/tree/master/subjects/bomberman-dom  

## dev recommendations:  
- do not use `zero` branch (it will be used finally)  
- create your branches from `dev` branch  
- name your own branches , f.e. `your-name-something` and work there  
- before make pull request to merge your changes to `dev` branch, create new branch `merge-to-dev-your-name-something`  
- create pull request to `merge-to-dev-your-name-something` from branch `your-name-something`  
- resolve any conflicts, and merge, now branch `merge-to-dev-your-name-something` ready to be merged to `dev`  
- create pull requests to `dev` branch, to merge your results into project  
- `dev` will be merged into `zero` branch, close to release  
---
- if you use more than one machine, press vscode `Synchronize Changes` button, every time you start new coding session, to decrease difficulties to push changes later. Extra conflict resolving can be extra headache, and always a potential way to make more mistakes. Even if it is your branch you can forget something accedentally.
