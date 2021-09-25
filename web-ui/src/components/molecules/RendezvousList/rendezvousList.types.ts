export interface IRendezvousListProps {
  trigger: boolean;

  setTrigger: (newState: boolean) => void;
}

export interface IRendezvousStatusResponse {
  address: string;
  status: string;
}
