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
    root = e.parentNode.parentNode.parentNode
    root.removeChild(e.parentNode.parentNode);
}

function getRoles() {
    let roles = [];
    let rolesRow = document.querySelectorAll('.role_row');

    let i = 0;
    rolesRow.forEach(function (row) {
        let role = new Role;
        let access = new Access;

        if (i == 0) {
            role.name = "User"
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
        role.id = i;
        i++;

        roles.push(role);
    })

    return roles;
}

function UpdateRoles() {
    let query = ""
    let newRoles = getRoles();

    for (let i = 0; i < newRoles.length; i++) {
        let toUpdate = false
        let roleA = newRoles[i];

        let j = 0;
        rolesWhenInit.forEach(function (roleB) {
            if (roleA.id == roleB.id) {
                if (roleA.name != roleB.name || roleA.access.create != roleB.access.create || roleA.access.pin != roleB.access.pin || roleA.access.remove != roleB.access.remove || roleA.access.ban != roleB.access.ban || roleA.access.manage != roleB.access.manage || roleA.access.give != roleB.access.give) {
                    toUpdate = true;
                    return
                }
            }

            j++;
        })

        if (toUpdate) {
            query += roleA.id + "," + roleA.name + "," + roleA.access.create + "," + roleA.access.pin + "," + roleA.access.remove + "," + roleA.access.ban + "," + roleA.access.manage + "," + roleA.access.give + ";";
        }
    }

    console.log(query);
}

let rolesWhenInit = getRoles()