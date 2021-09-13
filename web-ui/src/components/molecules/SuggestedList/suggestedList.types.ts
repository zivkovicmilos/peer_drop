export interface ISuggestedListProps {
  files: IFileInfo[];
}

export interface IFileInfo {
  id: string;
  name: string;
  extension: string;
}