const { default: fetch } = require("node-fetch")

module.exports = async function() {
  try {
    const res = await fetch('https://api.github.com/repos/garraflavatra/rolens');
    const json = await res.json();
    return { stars: json.stargazers_count || 0 };
  }
  catch (error) {
    console.error('Error: could not fetch GitHub stars:', error);
    return { stars: 0 };
  }
}
