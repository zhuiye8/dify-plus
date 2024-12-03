export default function randomPasswd(length) {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";

    function getRandomChar(set) {
        return set.charAt(Math.floor(Math.random() * set.length));
    }

    function generateRandomString(len) {
        let result = '';
        for (let i = 0; i < len; i++) {
            result += getRandomChar(charset);
        }
        return result;
    }

    function isValid(str) {
        return /[0-9]/.test(str) && /[a-zA-Z]/.test(str);
    }

    let password;
    do {
        password = generateRandomString(length);
    } while (!isValid(password));

    return password;
}