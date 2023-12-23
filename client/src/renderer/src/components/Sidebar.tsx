import ConversationList from './ConversationList'
import ConversationSearch from './ConversationSearch'

const Sidebar = () => {
  return (
    <div className="overflow-hidden h-full">
      <ConversationSearch />

      <ConversationList />
    </div>
  )
}

export default Sidebar
