import { ReportSharedStateVariable } from '$wails/go/app/App';
import { writable } from 'svelte/store';

function sharedStateStore(name) {
  const { set, subscribe } = writable();
  subscribe(newValue => ReportSharedStateVariable(name, newValue));
  return { set, subscribe };
}

const sharedState = {
  currentHost: sharedStateStore('currenthost'),
  currentDb: sharedStateStore('currentdb'),
  currentColl: sharedStateStore('currentcoll'),
};

export default sharedState;
