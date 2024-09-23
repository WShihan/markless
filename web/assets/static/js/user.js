function logout(route){
    document.cookie = `markless-token=; expires=-1;path=/`;
    window.location = route
}