import { StartProgressBar, StopProgressBar } from '$wails/go/ui/UI';

let taskCounter = 0;

export function startProgress(taskDescription = 'Loading…') {
  const taskIndex = ++taskCounter;
  StartProgressBar(taskIndex, taskDescription);

  const task = {
    id: taskIndex,
    description: taskDescription,
    end: () => StopProgressBar(taskIndex),
  };

  return task;
}
