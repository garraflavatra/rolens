import { StartProgressBar, StopProgressBar } from '$wails/go/ui/UI.js';

let taskCounter = 0;

export function startProgress(taskDescription = 'Loading…') {
  const taskIndex = ++taskCounter;
  let started = false;

  const debouncer = setTimeout(() => {
    StartProgressBar(taskIndex, taskDescription);
    started = true;
  }, 150);

  const task = {
    id: taskIndex,
    description: taskDescription,

    end: () => {
      clearTimeout(debouncer);
      if (started) {
        StopProgressBar(taskIndex);
      }
    },
  };

  return task;
}
