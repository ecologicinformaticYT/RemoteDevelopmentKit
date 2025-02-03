//global variables

let server = '';
let password = '';
let mode = '';

let proj = "";
let path="";
let fileopened = false;

let lastXHRstatus = -1;

let autoDdelay = undefined //automatic refresh for slow/dual debug

let fcmd = "" //fast command
let scmd = "" //slow command
let dbgId = "" //debug session ID

//functions
async function xhrPOST(route, data) {
    /*
    Sends a XHR POST request to server/route
    Parameters :
        (string) route : the route
        (Object) data : dictionary made with the data to communicate to the server
    */

    return new Promise((resolve, reject) => {
        // Create a new instance of XMLHttpRequest
        var xhr = new XMLHttpRequest();

        // Configure the request
        xhr.open("POST", `${server}/${route}`, true);
        xhr.setRequestHeader("Content-Type", "application/json");

        // Define a callback function to handle the response
        xhr.onreadystatechange = () => {
            if (xhr.readyState === 4) {
                if (xhr.status === 200 || xhr.status === 201) {
                    lastXHRstatus = 200;
                    resolve(xhr.responseText);  //resolve the promise

                } else {
                    lastXHRstatus = xhr.status;
                    alert(`XHR POST failed: STATUS/CODE: ${xhr.status} | STEXT: ${xhr.statusText} | RTEX: ${xhr.responseText}`);
                    reject(new Error(xhr.responseText)); // Rejeter la promesse en cas d'erreur
                }
            }
        };

        // Send the request with the JSON data
        xhr.send(JSON.stringify(data));
    });
}

async function getProjectsList(){
    /*
    Get the list of projects stored on the server
    */
    let res = await xhrPOST("listProjects", {"mode":`${mode}`,"password":`${password}`})
    
    if (lastXHRstatus == 200 && res != '' && res != undefined){
        let projs = String(res).split("\n")

        projs.forEach(proj_ => {
            if (proj_ != "" && proj_ != "\n"){
                let option = document.createElement("option")
                option.id = `${proj_}`
                option.textContent = `${proj_}`

                document.getElementById('project').appendChild(option)
            }
        });
    
    } 
}

async function login(){
    /*
    Login to the server
    */
    server=document.getElementById('server').value
    server+=`:${document.getElementById('port').value}`

    password=document.getElementById('password').value

    let res = await xhrPOST("login", {"mode":`${mode}`,"password":`${password}`})
    

    if (lastXHRstatus == 200 && res != ''){
        document.getElementById('settings').style.visibility = 'hidden'
        document.getElementById('settings').style.height = '0px'

        document.getElementById('lgout').style.visibility = 'visible'
        document.getElementById('projSelector').style.visibility = 'visible'

        document.getElementById('wassup').style.visibility = 'hidden'

        getProjectsList()
    }
}

function npWindow(state){
    /*
    Display or stop displaying the create project div element
    */
    if (state == "open"){
        document.getElementById('cp').style.height = "10%"
        document.getElementById('cp').style.visibility = "visible"
    } else if (state == "close"){
        document.getElementById('cp').style.height = "0%"
        document.getElementById('cp').style.visibility = "hidden"
    }
}

async function openProject(){
    /*
    Loads a project architecture
    */
    if (document.getElementById('cp').style.visibility == 'visible'){
        npWindow('close')
    }

    proj = document.getElementById('project').value;

    if (proj == "" || proj == undefined){
        return
    }

    let res = await xhrPOST("listProjectBranches", {"mode":mode,"password":password,"project":proj})

    if (lastXHRstatus == 200){
        document.getElementById('IDE').style.visibility = 'visible'
        document.getElementById('debugger').style.visibility = 'visible'

        let branches = res.split('\n')
        branches.forEach(branch => {
            if (branch != "" && branch != "\n"){
                let option = document.createElement("option")
                option.id = `${branch}`
                option.textContent = `${branch}`

                document.getElementById('bhelper').appendChild(option)
            }
        });

        document.getElementById('projSelector').style.visibility = 'hidden'
        document.getElementById('projSelector').style.height = '0px'
    }
}

async function delProj(){
    /*
    Delete (forever) a project on the server
    */
    proj = document.getElementById('project').value;

    if (proj == ''){
        return
    }

    mode='admin'

    let res = await xhrPOST("delete", {"mode":mode,"password":password,"project":proj,"path":""})
    if (lastXHRstatus == 200 && res != ''){
        alert(`Project: ${proj} was deleted sucessfuly !`)
        history.go(0)
    }

    mode=''
}

async function np(){
    proj = document.getElementById('projectName').value

    mode='admin'

    let res = await xhrPOST("mkdir", {"mode":mode,"password":password,"project":proj,"path":""})
    if (lastXHRstatus == 200 && res != ''){
        alert(`Project: ${proj} was created sucessfuly !`)
    }

    npWindow('close')
    getProjectsList()

    mode='' 
}

async function closeFile(){
    /*
    Unlock a file that was locked for edition
    */
    if (path == '' || proj == ''){
        return
    }

    let res = await xhrPOST("close", {"mode":mode,"password":password,"project":proj,"path":path})
    if (lastXHRstatus == 200 && res != ''){
        alert('File closed sucessfuly !')
        editor.setValue('')
        path=''
        fileopened = false
    }
}

async function openFile(){
    /*
    Opens (and locks for edition) a file of the project
    */
    if (fileopened == true){
        await closeFile() //close the previously opened file 
    }

    path = document.getElementById('thePath').value //get the path of the file to open

    if (path == '' && (document.getElementById('bhelper').value != '' && document.getElementById('bhelper').value != undefined)){
        path = document.getElementById('bhelper').value
        document.getElementById('thePath').value = path
    } else if (path == '') {
        alert("Please, select a file to open.")
        return
    }

    let res = await xhrPOST("read", {"mode":mode,"password":password,"project":proj,"path":path})
    if (lastXHRstatus == 200){
        editor.setValue(res)
        fileopened = true
    }

}

async function deleteFile(){
    /*
    Sends a request to delete a file
    */

    mode='admin'

    path = document.getElementById('thePath').value

    closeFile()

    if (path == ''){
        path = document.getElementById('bhelper').value
        document.getElementById('thePath').value = path
    } else if (path == ''){
        return
    }
    
    let res = await xhrPOST("delete", {"mode":mode,"password":password,"project":proj,"path":path})

    if (lastXHRstatus == 200 && res != ''){
        openProject()
        alert(`${path} (file) was deleted sucessfuly !`)
    }

    mode=''
}

async function createFile(){
    /*
    Creates a new file on the server
    */
    path = document.getElementById('thePath').value

    if (path == ''){
        return
    } 

    let content = editor.getValue()

    let res = await xhrPOST("writef", {"mode":mode,"password":password,"project":proj,"path":path,"content":content})
 
    if (lastXHRstatus == 200 && res != ''){
        alert(`${path} saved sucessfuly !`)
        openProject()
        document.getElementById('save_status').value = `Status : Saved !`;
    }
}

async function makedir(){
    /*
    Creates a new directory
    */
    path = document.getElementById('thePath').value

    if (path == ''){
        path = document.getElementById('bhelper').value
        document.getElementById('thePath').value = path
    } else if (path == ''){
        return
    }

    let res = await xhrPOST("mkdir", {"mode":mode,"password":password,"project":proj,"path":path})
 
    if (lastXHRstatus == 200 && res != ''){
        let branches = res.split('\n')
        branches.forEach(branch => {
            document.getElementById('bhelper').appendChild(`<option id=${branch}>${branch}</option>`) 
        });
        alert('Directory created sucessfuly !')
    }
}

async function save(){
    /*
    Saves the new content of the file on the server
    */

    path = document.getElementById('thePath').value

    if (path == ''){
        path = document.getElementById('bhelper').value
        document.getElementById('thePath').value = path
    } else if (path == '') {
        return
    }
    let content = editor.getValue()

    let res = await xhrPOST("write", {"mode":mode,"password":password,"project":proj,"path":path,"content":content})
 
    if (lastXHRstatus == 200 && res != ''){
        alert(`${path} saved sucessfuly !`)
    }
}

async function rmv(){
    /*
    Deletes a directory on the server
    */
    mode='admin'
    
    path = document.getElementById('thePath').value

    if (path == ''){
        path = document.getElementById('bhelper').value
        document.getElementById('thePath').value = path
    } else if (path == ''){
        return
    }
    
    let res = await xhrPOST("delete", {"mode":mode,"password":password,"project":proj,"path":path})
 
    if (lastXHRstatus == 200 && res != ''){
        alert(`${path} deleted sucessfuly !`)
    }

    mode=''
}

async function fastDebug(){
    /*
    Uses remote-debug (fast way)
    */
    fcmd = document.getElementById('dbgIO').value

    if (fcmd == ''){
        return
    }

    let res = await xhrPOST("fdebug", { "mode": mode, "password": password, "project": proj, "path": path, "cmd":fcmd})
 
    if (lastXHRstatus == 200 && res != ''){
        document.getElementById('dbgIO').value = `(Fast) command:\n${fcmd}\n\nRESULT:\n${res}`
    }

    mode=''
}

function delDelay(dict){
    /*
    Stops the recurrent call to debugRecall when the slow debug results come
    */
    if (dict['slow_output'] != ""){
        autoDdelay=undefined
    }
    
}

async function debugRecall(){
    /*
    Retrieves slow/dual debug results
    */
    let res = await xhrPOST("debugRecall", { "mode": mode, "password": password, "session": dbgId})

    if (lastXHRstatus == 200 && res != ''){
        let res_ = createDictFromJSON(res)
        delDelay(res_)

        document.getElementById('dbgIO').value = `(Fast) command:\n${fcmd}\n\nRESULT:\n${res_['output']}`
        document.getElementById('sdbgIO').value = `(Slow) command:\n${fcmd}\n\nRESULT:\n${res_['slow_output']}`
    }
}

async function slowDebug(){
    /*
    Start a slow debugging session
    */
    scmd = document.getElementById('sdbgIO').value

    if (scmd == ''){
        return
    }

    let res = await xhrPOST("debug", { "mode": mode, "password": password, "project": proj, "path": path,"slow_cmd":scmd, "cmd":""})
 
    if (lastXHRstatus == 200 && res != ''){
        dbgId=res
        autoDdelay=setTimeout(debugRecall, 2000)
    }

    mode=''
}

async function dualDebug(){
    /*
    Starts a dual debuggin session
    */
    fcmd = document.getElementById('dbgIO').value
    scmd = document.getElementById('sdbgIO').value

    if (scmd == '' || fcmd == ''){
        return
    }

    let res = await xhrPOST("debug", { "mode": mode, "password": password, "project": proj, "path": path,"slow_cmd":scmd, "cmd":""})
 
    if (lastXHRstatus == 200 && res != ''){
        dbgId=res
        autoDdelay=setTimeout(debugRecall, 2000)
    }

    mode=''
}

function createDictFromJSON(jsonString) {
    try {
        // Parse le JSON en objet JavaScript
        let parsed = JSON.parse(jsonString);

        // Si c'est un tableau, on prend le premier élément (en supposant qu'il contienne des objets)
        if (Array.isArray(parsed)) {
            if (parsed.length > 0 && typeof parsed[0] === 'object') {
                return parsed[0]; // On retourne le premier objet du tableau
            } else {
                throw new Error("the JSON table doesn't have any valid object.");
            }
        } else if (typeof parsed === 'object') {
            return parsed; // Si c'est déjà un objet, on le retourne directement
        } else {
            throw new Error("the JSON doesn't include any element/table of elements.");
        }
    } catch (error) {
        console.error("Error while parsing JSON:", error);
        return null;
    }
}

function setPath(){
    /*
    Set the path variable to the content of the related html element
    */
   path=document.getElementById('thePath').value
   alert(`Path variable = ${path}`)
}