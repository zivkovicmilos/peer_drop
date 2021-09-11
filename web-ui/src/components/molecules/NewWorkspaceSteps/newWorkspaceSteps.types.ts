import { ENewWorkspaceStep } from '../../pages/NewWorkspace/newWorkspace.types';

export interface INewWorkspaceStepsProps {
  currentStep: ENewWorkspaceStep;
  currentStepIndex: number;
  steps: ENewWorkspaceStep[];
}
