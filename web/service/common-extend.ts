import { get } from '@/service/base'
import type { UserMoney } from '@/models/common-extend'

export const fetchUserMoney = () => {
  return get<UserMoney>('account/money')
}
