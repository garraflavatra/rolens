#!/usr/bin/env node

const { execSync, spawn } = require('child_process');
const { readFileSync, statSync, rmdirSync, mkdirSync } = require('fs');

// Check that the script is run from the root.

try {
  const wailsJsonFile = statSync('./wails.json');
  if (!wailsJsonFile.isFile()) {
    throw new Error();
  }
} catch {
  console.log('Error: please run the build script from the Rolens project root.');
  process.exit(1);
}

// Output version.

const version = JSON.parse(readFileSync('./wails.json').toString()).info.productVersion;

if (process.argv.includes('-v') || process.argv.includes('-version') || process.argv.includes('--version')) {
  console.log(version);
  process.exit(0);
}

// Output help text.

if (process.argv.includes('-h') || process.argv.includes('--help')) {
  console.log(`Rolens build script v${version}`);
  console.log('');
  console.log('This script installs missing dependencies if any, and then compiles Rolens');
  console.log('for the current platform.');
  console.log('');
  console.log('Options:');
  console.log('  -h --help    Show this help text and exit.');
  console.log('  -q --quiet   Do not output Wails build log.');
  console.log('  -v --version Log the current Rolens version and exit.');

  process.exit(0);
}

// Shared objects.

const quiet = process.argv.includes('-q') || process.argv.includes('--quiet');
const isWindows = process.platform === 'win32';
const missingDependencies = [];

function isNullish(val) {
  return val === undefined || val === null;
}

// Check that Go ^1.20 is installed.

try {
  const goMinorVersion = /go1\.([0-9][0-9])/.exec(
    execSync('go version').toString()
  )?.pop();

  if (isNullish(goMinorVersion) || (parseInt(goMinorVersion) < 20)) {
    throw new Error();
  }
} catch {
  missingDependencies.push({ name: 'Go ^1.20', url: 'https://go.dev/doc/install' });
}

// Check that Node.js ^16 is installed.

try {
  const nodeMajorVersion = /v([0-9]{1,2})\.[0-9]{1,3}\.[0-9]{1,3}/.exec(
    execSync('node --version').toString()
  )?.pop();

  if (isNullish(nodeMajorVersion) || (parseInt(nodeMajorVersion) < 16)) {
    throw new Error();
  }
} catch {
  missingDependencies.push({ name: 'Node.js ^16', url: 'https://go.dev/doc/install' });
}

// Check that Wails is installed.

try {
  const wailsMinorVersion = /v2\.([0-9])\.[0-9]/.exec(
    execSync('wails version').toString()
  )?.pop();

  if (isNullish(wailsMinorVersion) || (parseInt(wailsMinorVersion) < 3)) {
    throw new Error();
  }
} catch {
  missingDependencies.push({
    name: 'Wails ^2.3',
    command: 'go install github.com/wailsapp/wails/v2/cmd/wails@latest',
    url: 'https://wails.io/docs/gettingstarted/installation',
  });
}

// Check that NSIS is installed on Windows.

if (isWindows) {
  try {
    const nsisInstalled = /v3\.([0-9][0-9])/.test(execSync('makensis.exe /VERSION').toString());
    if (!nsisInstalled) {
      throw new Error();
    }
  } catch {
    missingDependencies.push({
      name: 'Nullsoft Install System ^3',
      command: 'choco install nsis',
      url: 'https://nsis.sourceforge.io/Download',
      comment: 'Note: you should add makensis.exe to your path:\n    setx /M PATH "%PATH%;C:\\Program Files (x86)\\NSIS\\Bin"'
    });
  }
}

// Report missing dependencies.

if (missingDependencies.length > 0) {
  console.log('You are missing the following dependencies:');

  for (const dependency of missingDependencies) {
    console.log('');
    console.log(`- ${dependency.name}`);

    if (dependency.command) {
      console.log('  Install it by executing:');
      console.log(`    ${dependency.command}`);
    }

    if (dependency.url) {
      console.log('  Visit the following page for more information:');
      console.log(`    ${dependency.url}`);
    }

    if (dependency.comment) {
      console.log(`  ${dependency.comment}`);
    }
  }

  process.exit(1);
}

// Clean output directory.

console.log('Cleaning output directory...');
try { rmdirSync('./build/bin'); } catch {}
try {
  mkdirSync('./build/bin');
}
catch (err) {
  console.log('Failed to create build output directory!');
}

// Build Rolens.

console.log(`Building Rolens ${version}...`);
console.log();

const proc = spawn('wails', [ 'build', '-clean', isWindows ? '-nsis' : '' ]);

if (!quiet) {
  const suppressMessages = [
    'Wails CLI',
    'If Wails is useful',
    'https://github.com/sponsors/leaanthony',
  ];

  proc.stdout.on('data', data => {
    for (let i = 0; i < suppressMessages.length; i++) {
      if (data.toString().indexOf(suppressMessages[i]) !== -1) {
        return;
      }
    }
    process.stdout.write(data);
  });

  proc.stderr.on('data', data => process.stderr.write(data));
}

proc.on('exit', code => {
  console.log();
  process.exit(code);
});
