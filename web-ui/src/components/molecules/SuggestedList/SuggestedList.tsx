import { Box } from '@material-ui/core';
import { FC } from 'react';
import SuggestedFile from '../../atoms/SuggestedFile/SuggestedFile';
import { ISuggestedListProps } from './suggestedList.types';

const SuggestedList: FC<ISuggestedListProps> = (props) => {
  const { files, workspaceMnemonic } = props;

  return (
    <Box display={'flex'} width={'80%'} ml={-4}>
      {files.map((file) => {
        return (
          <SuggestedFile file={file} workspaceMnemonic={workspaceMnemonic} />
        );
      })}
    </Box>
  );
};

export default SuggestedList;