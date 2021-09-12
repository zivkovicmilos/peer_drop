import { ENWPermissionsType } from '../../molecules/NewWorkspacePermissions/newWorkspacePermissions.types';

export interface ISwitchBoxProps {
  type: ENWPermissionsType;
  description: string;
  onToggle: (state: boolean) => void;
  toggled: boolean;

  disabled: boolean;
}
