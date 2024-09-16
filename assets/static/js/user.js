function logout(route){
    document.cookie = `markee-token=; expires=-1;path=/`;
    console.log('登出')
    window.location = route
}