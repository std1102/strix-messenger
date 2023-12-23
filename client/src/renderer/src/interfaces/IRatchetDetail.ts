interface IRatchetDetail {
  root_key: string
  chain_send_key: string
  chain_recv_key: string
  total_message_sent: number
  total_message_recv: number
  missing_message_keys: null | string[]
  rachet_id: string
}

export default IRatchetDetail
