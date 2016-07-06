// Checks if string 1 is a fuzzy match of string 2 (if string1 is in string2)
module.exports = function (str1, str2) {
  str1 = str1.toLowerCase().replace(/[\W_]/g, '');
  str2 = str2.toLowerCase().replace(/[\W_]/g, '');

  if (str2.indexOf(str1) > -1) {
    return true;
  }

  return false;
}
