import { Box } from '@material-ui/core';
import AddRoundedIcon from '@material-ui/icons/AddRounded';
import { FC, useEffect, useState } from 'react';
import { useHistory } from 'react-router-dom';
import IdentitiesService from '../../../services/identities/identitiesService';
import { IIdentityResponse } from '../../../services/identities/identitiesService.types';
import ActionButton from '../../atoms/ActionButton/ActionButton';
import PageTitle from '../../atoms/PageTitle/PageTitle';
import Pagination from '../../atoms/Pagination/Pagination';
import usePagination from '../../atoms/Pagination/pagination.hook';
import IdentityCard from '../../molecules/IdentityCard/IdentityCard';
import IdentitySort from '../../molecules/IdentitySort/IdentitySort';
import {
  EIdentitySortDirection,
  EIdentitySortParam
} from '../../molecules/IdentitySort/identitySort.types';
import useSnackbar from '../../molecules/Snackbar/useSnackbar.hook';
import { IIdentitiesProps } from './identities.types';

const Identities: FC<IIdentitiesProps> = () => {
  const history = useHistory();

  const [activeSort, setActiveSort] = useState<EIdentitySortParam>(
    EIdentitySortParam.NAME
  );
  const [sortDirection, setSortDirection] = useState<EIdentitySortDirection>(
    EIdentitySortDirection.ASC
  );

  const [identities, setIdentities] = useState<{ data: IIdentityResponse[] }>({
    data: []
  });

  const handleNewIdentity = () => {
    history.push('/identities/new');
  };

  const { page, count, setCount, limit, handlePageChange } = usePagination({
    limit: 8
  });

  const { openSnackbar } = useSnackbar();

  useEffect(() => {
    const fetchIdentities = async () => {
      return await IdentitiesService.getIdentities({ page, limit });
    };

    fetchIdentities()
      .then((response) => {
        setIdentities({ data: response.data });
        setCount(response.count);
      })
      .catch((err) => {
        openSnackbar('Unable to fetch identities', 'error');

        setCount(0);
      });
  }, [page]);

  // TODO add useEffect listeners for the sort params / direction

  return (
    <Box
      display={'flex'}
      flexDirection={'column'}
      width={'100%'}
      height={'100%'}
    >
      <Box
        display={'flex'}
        justifyContent={'space-between'}
        width={'100%'}
        alignItems={'center'}
        mb={4}
      >
        <PageTitle title={'Identities'} />
        <ActionButton
          text={'New identity'}
          startIcon={<AddRoundedIcon />}
          onClick={handleNewIdentity}
        />
      </Box>

      <Box display={'flex'} width={'100%'} mb={4}>
        <IdentitySort
          activeSort={activeSort}
          setActiveSort={setActiveSort}
          setSortDirection={setSortDirection}
          sortDirection={sortDirection}
        />
      </Box>

      <Box display={'flex'} width={'100%'} flexDirection={'column'}>
        <Box
          display={'flex'}
          flexWrap={'wrap'}
          width={'100%'}
          ml={'-36px'}
          mt={'-36px'}
        >
          {identities.data.map((identity) => {
            return (
              <IdentityCard
                key={`identity-${identity.id}`}
                id={identity.id}
                isPrimary={identity.isPrimary}
                picture={identity.picture}
                name={identity.name}
                publicKeyID={identity.publicKeyID}
                numWorkspaces={identity.numWorkspaces}
                dateCreated={identity.dateCreated}
              />
            );
          })}
        </Box>
        <Pagination
          count={count}
          limit={limit}
          page={page}
          onPageChange={handlePageChange}
        />
      </Box>
    </Box>
  );
};

export default Identities;
