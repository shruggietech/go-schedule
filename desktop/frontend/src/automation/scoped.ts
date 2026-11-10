import type { AutomationBridge } from './model'
// Every mutation carries the identity rendered when its control was opened.
export function scopedAutomation(bridge: AutomationBridge, daemonID: string): AutomationBridge {
 return { ...bridge,
  workspace: () => bridge.workspace(daemonID),
  saveChain: d => bridge.saveChain({...d, expectedDaemonId:daemonID}), deleteChain:id=>bridge.deleteChain(id,daemonID),
  saveTrigger:d=>bridge.saveTrigger({...d,expectedDaemonId:daemonID}), setTriggerEnabled:(id,v)=>bridge.setTriggerEnabled(id,v,daemonID), revealTrigger:id=>bridge.revealTrigger(id,daemonID), rotateTrigger:id=>bridge.rotateTrigger(id,daemonID), fireTrigger:id=>bridge.fireTrigger(id,daemonID), deleteTrigger:id=>bridge.deleteTrigger(id,daemonID),
  createTriggerSet:d=>bridge.createTriggerSet({...d,expectedDaemonId:daemonID}), retargetTriggerSet:(id,target,original,overwrite)=>bridge.retargetTriggerSet(id,target,original,overwrite,daemonID), setTriggerSetEnabled:(id,v)=>bridge.setTriggerSetEnabled(id,v,daemonID), revealTriggerSet:id=>bridge.revealTriggerSet(id,daemonID),rotateTriggerSet:id=>bridge.rotateTriggerSet(id,daemonID),deleteTriggerSet:id=>bridge.deleteTriggerSet(id,daemonID),
  saveWatcher:d=>bridge.saveWatcher({...d,expectedDaemonId:daemonID}),setWatcherEnabled:(id,v)=>bridge.setWatcherEnabled(id,v,daemonID),deleteWatcher:id=>bridge.deleteWatcher(id,daemonID),
  joinDefinition:bridge.joinDefinition ? id=>bridge.joinDefinition!(id,daemonID):undefined,
  saveJoin:bridge.saveJoin ? d=>bridge.saveJoin!({...d,expectedDaemonId:daemonID}):undefined,deleteJoin:bridge.deleteJoin ? id=>bridge.deleteJoin!(id,daemonID):undefined,joinInstances:bridge.joinInstances ? (id,offset)=>bridge.joinInstances!(id,offset,daemonID):undefined,joinInstance:bridge.joinInstance ? id=>bridge.joinInstance!(id,daemonID):undefined,
 }
}
