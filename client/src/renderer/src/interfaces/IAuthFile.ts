interface IIdentityKey {
  public_key: string
  private_key: string
  public_hash: string
  private_hash: string
}

interface IAuthFile {
  identity_key: IIdentityKey
  one_time_key: {
    [key: string]: IIdentityKey
  }
  pre_key: IIdentityKey
  pin: string
}

export default IAuthFile
