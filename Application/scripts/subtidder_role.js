const table = document.querySelector('#admin_roles');

class Role {
    constructor(id, name, access) {
        this.id = id;
        this.name = name;
        this.access = access;
    }
}

class Access {
    constructor(id, create, pin, remove, ban, manage, give) {
        this.id = id;
        this.create = create;
        this.pin = pin;
        this.remove = remove;
        this.ban = ban;
        this.manage = manage;
        this.give = give;
    }
}

function AddRoleRow() {
    CreateRow();
}

function CreateRow() {
    let tr = document.createElement('tr');
    tr.classList.add('role_row');
    tr.id = "role_NEW_" + Math.floor(Math.random() * 10000000);

    let nameCol = document.createElement('td');
    let nameInput = document.createElement('input');
    nameInput.setAttribute('type', 'text');
    nameInput.setAttribute('placeholder', 'Role Name');
    nameCol.appendChild(nameInput);

    let createCol = document.createElement('td');
    let createInput = document.createElement('input');
    createInput.setAttribute('type', 'checkbox');
    createCol.appendChild(createInput);

    let pinCol = document.createElement('td');
    let pinInput = document.createElement('input');
    pinInput.setAttribute('type', 'checkbox');
    pinCol.appendChild(pinInput);

    let removeCol = document.createElement('td');
    let removeInput = document.createElement('input');
    removeInput.setAttribute('type', 'checkbox');
    removeCol.appendChild(removeInput);

    let banCol = document.createElement('td');
    let banInput = document.createElement('input');
    banInput.setAttribute('type', 'checkbox');
    banCol.appendChild(banInput);

    let manageCol = document.createElement('td');
    let manageInput = document.createElement('input');
    manageInput.setAttribute('type', 'checkbox');
    manageCol.appendChild(manageInput);

    let giveCol = document.createElement('td');
    let giveInput = document.createElement('input');
    giveInput.setAttribute('type', 'checkbox');
    giveCol.appendChild(giveInput);

    let removeCol2 = document.createElement('td');
    let p = document.createElement('p');
    p.innerHTML = "X";
    p.style.color = "red";
    p.classList.add('remove_row');
    p.setAttribute("onclick", "removeRow(this)");
    removeCol2.appendChild(p);

    tr.appendChild(nameCol);
    tr.appendChild(createCol);
    tr.appendChild(pinCol);
    tr.appendChild(removeCol);
    tr.appendChild(banCol);
    tr.appendChild(manageCol);
    tr.appendChild(giveCol);
    tr.appendChild(removeCol2);

    insertBefore(tr, table.querySelector('#create_role'));
}

function insertBefore(newNode, existingNode) {
    existingNode.parentNode.insertBefore(newNode, existingNode);
}

function removeRow(e) {
    root = e.parentNode.parentNode
    root.removeChild(e.parentNode);
}

function getRoles() {
    let roles = [];
    let rolesRow = document.querySelectorAll('.role_row');
    

    let i = 0;
    rolesRow.forEach(function (row) {
        let role = new Role;
        let access = new Access;

        let id = "NEW";
        if (!row.id.includes("NEW")) {
            id = row.id.replace("role_", "");
        }

        if (id == "-1") {
            role.name = "User"
            role.access = access;
            role.id = id;
            i++;
            roles.push(role);
            return;
        } else {
            role.name = row.querySelector(':nth-child(1)').firstChild.value;
        }
        
        access.create = row.querySelector(':nth-child(2)').firstChild.checked;
        access.pin = row.querySelector(':nth-child(3)').firstChild.checked;
        access.remove = row.querySelector(':nth-child(4)').firstChild.checked;
        access.ban = row.querySelector(':nth-child(5)').firstChild.checked;
        access.manage = row.querySelector(':nth-child(6)').firstChild.checked;
        access.give = row.querySelector(':nth-child(7)').firstChild.checked;

        role.access = access;
        role.id = id;
        i++;

        roles.push(role);
    })

    return roles;
}

let rolesWhenInit = getRoles()
function UpdateRoles() {
    let toCreate = "";
    let toChange = "";
    let toDelete = "";

    let newRoles = getRoles();

    let toCreateArr = [];
    let toDeleteArr = [];

    console.log("Before\nrolesWhenInit : ");
    printArr(rolesWhenInit);
    console.log("newRoles : ");
    printArr(newRoles);

    for (let i = 0; i < rolesWhenInit.length; i++) {
        let role = rolesWhenInit[i];
        let stillExist = false;
        
        if (role == null) {continue;}
        for (let j = 0; j < newRoles.length; j++) {
            let newRole = newRoles[j];
            if (newRole == null) {continue;}
            if (role == null) {continue;}
            if (role.id == newRole.id) {
                stillExist = true;
                if (!compareRoles(role, newRole)) { // toChange
                    toChange += newRole.id + "," + newRole.name + "," + newRole.access.create + "," + newRole.access.pin + "," + newRole.access.remove + "," + newRole.access.ban + "," + newRole.access.manage + "," + newRole.access.give + ";";
                }
                delete newRoles[j];
                delete rolesWhenInit[i];
            }
        }

        if (!stillExist) {
            toDelete += role.id + ";";
            delete rolesWhenInit[i];
        }
    }

    console.log("After\nrolesWhenInit : ");
    printArr(rolesWhenInit);
    console.log("newRoles : ");
    printArr(newRoles);


    newRoles.forEach(function (newRole) {   // toCreate
        if (newRole == null) {return;}
        toCreate += newRole.id + "," + newRole.name + "," + newRole.access.create + "," + newRole.access.pin + "," + newRole.access.remove + "," + newRole.access.ban + "," + newRole.access.manage + "," + newRole.access.give + ";"; 
    })

    console.log("\ntoChange :", toChange);
    console.log("toCreate :", toCreate);
    console.log("toDelete :", toDelete);
    
    fetch(location.pathname, {
        method: "post",
        headers: {
          'Content-Type': 'application/json'
        },
      
        //make sure to serialize your JSON body
        body: JSON.stringify({
          "role_to_create": toCreate,
          "role_to_update": toChange,
          "role_to_delete": toDelete,
        })
    }).then(() => {
        window.location.reload();
    })
    
}

function compareRoles(roleA, roleB) {
    if (roleA.id != roleB.id) {
        return false;
    } else {
        if (roleA.name != roleB.name || roleA.access.create != roleB.access.create || roleA.access.pin != roleB.access.pin || roleA.access.remove != roleB.access.remove || roleA.access.ban != roleB.access.ban || roleA.access.manage != roleB.access.manage || roleA.access.give != roleB.access.give) {
            return false;
        }
    }
    return true;
}

function printArr(arr) {
    for (let i = 0; i < arr.length; i++) {
        if (arr[i] == null) {
            console.log("null");
        } else {
            console.log(arr[i]);       
        }
    }
}
