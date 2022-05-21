let atributionWhenInit = getAtribution();
console.log("atributionWhenInit : ", atributionWhenInit);

function getAtribution() {
    let atribution = [];
    
    document.querySelectorAll(".role_attribution_row").forEach(function(e) {
        let select = e.querySelector('select');
        let roleId = "-69420";      // -69420 is owner role id; -1 is user role id
        if (select) {
            roleId = select.options[select.selectedIndex].value;
        }
        let userId = e.querySelector('a').getAttribute("user_id");

        atribution.push({
            userId: userId,
            roleId: roleId
        })
    });

    return atribution;
}


function UpdateRoleAtribution() {
    let newAtribution = getAtribution();

    let query = "";

    for (let i = 0; i < atributionWhenInit.length; i++) {
        let atr = atributionWhenInit[i];
        for (let j = 0; j < newAtribution.length; j++) {
            let newAtr = newAtribution[j];
            if (atr.userId == newAtr.userId && atr.roleId != newAtr.roleId) {
                query += newAtr.userId + "," + newAtr.roleId + ";";
            }
        }
    }

    console.log("query : ", query);
    return query;
}

function UpdateUser() {
    let banUpdate = UpdateBannedUser();
    let atributionUpdate = UpdateRoleAtribution();
  
    fetch(location.pathname, {
      method: "post",
      headers: {
        'Content-Type': 'application/json'
      },
    
      //make sure to serialize your JSON body
      body: JSON.stringify({
        "banned_user_changes": banUpdate,
        "role_atribution_changes": atributionUpdate
      })
    }).then(() => {
      window.location.reload();
    })
  }