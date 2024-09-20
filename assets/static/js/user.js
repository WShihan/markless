function logout(route){
    document.cookie = `markless-token=; expires=-1;path=/`;
    console.log('登出')
    window.location = route
}