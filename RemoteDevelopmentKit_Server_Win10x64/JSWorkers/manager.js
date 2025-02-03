const { listAllChildsSync, listChilds, write, read } = require('./modules/fsFacilities.js');

//global variables
let projects = []
let childs = []

let prevlen = parseInt(read('../cache/manager_cache.txt'))

//length function
function len(li){
    let i = 0
    li.forEach(element => {
        i++       
    });
    return i
}

console.log("RemoteDevelopmentKit JS Files Manager is now running...")

//mainloop
while (1+1==2){

    projects = listChilds("../projects/")

    if (len(projects) != prevlen || read("../cache/mc2.txt") == "1"){

        write("../projects/__list__.txt", "")
        projects.forEach(proj => {

            if (proj != "__list__.txt"){

                write("../projects/__list__.txt", `${read("../projects/__list__.txt")}${proj}\n`)
        
                childs = listAllChildsSync(`../projects/${proj}/`)
                write(`../projects/${proj}/__architecture__.txt`, "")
                childs.forEach(c => {
                    if (c.replace(`..\\projects\\${proj}\\`, "") != `__architecture__.txt` && !c.includes("locker")){
                        write(`../projects/${proj}/__architecture__.txt`, `${read(`../projects/${proj}/__architecture__.txt`)}${c.replace(`..\\projects\\${proj}\\`, "")}\n`)
                    }
                });
            }
        });
        prevlen = len(projects)
        write("../cache/mc2.txt", "0")
        write('../cache/manager_cache.txt', `${prevlen.toString()}`)
    } 
}