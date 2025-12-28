import { useState, useEffect } from 'react'
import { Plus, Trash2, Pencil, X, Save, RefreshCw } from 'lucide-react'

function App() {
  const [items, setItems] = useState([])
  const [loading, setLoading] = useState(true)
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [currentItem, setCurrentItem] = useState(null)
  
  // API URL - use proxy in dev, relative in prod
  const API_URL = import.meta.env.VITE_API_URL || '/api'

  const fetchItems = async () => {
    setLoading(true)
    try {
      const res = await fetch(`${API_URL}/items`)
      if (res.ok) {
        const data = await res.json()
        setItems(data || [])
      }
    } catch (error) {
      console.error("Failed to fetch items", error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchItems()
  }, [])

  const handleDelete = async (id) => {
    if (!confirm('Are you sure you want to delete this item?')) return
    try {
      await fetch(`${API_URL}/items/${id}`, { method: 'DELETE' })
      setItems(items.filter(item => item.id !== id))
    } catch (error) {
      console.error("Failed to delete", error)
    }
  }

  const handleSave = async (e) => {
    e.preventDefault()
    const formData = new FormData(e.target)
    const payload = {
      name: formData.get('name'),
      price: parseFloat(formData.get('price'))
    }

    try {
      if (currentItem) {
        await fetch(`${API_URL}/items/${currentItem.id}`, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload)
        })
      } else {
        await fetch(`${API_URL}/items`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload)
        })
      }
      setIsModalOpen(false)
      setCurrentItem(null)
      fetchItems()
    } catch (error) {
      console.error("Failed to save", error)
    }
  }

  const openModal = (item = null) => {
    setCurrentItem(item)
    setIsModalOpen(true)
  }

  return (
    <div className="app-layout">
      <div className="container">
        <header className="header">
          <div>
            <h1 className="title">
              Item Manager
            </h1>
            <p className="subtitle">Manage your inventory with style</p>
          </div>
          <button 
            onClick={() => openModal()}
            className="btn btn-primary"
          >
            <Plus size={20} />
            Add New Item
          </button>
        </header>

        {loading ? (
          <div className="loader-container">
            <RefreshCw size={32} className="spinner" />
          </div>
        ) : (
          <div className="items-grid">
            {items.map(item => (
              <div key={item.id} className="glass-panel item-card">
                <div className="item-header">
                  <div className="item-icon">
                    {item.name.charAt(0).toUpperCase()}
                  </div>
                  <div className="item-actions">
                    <button 
                      onClick={() => openModal(item)}
                      className="icon-btn"
                    >
                      <Pencil size={18} />
                    </button>
                    <button 
                      onClick={() => handleDelete(item.id)}
                      className="icon-btn danger"
                    >
                      <Trash2 size={18} />
                    </button>
                  </div>
                </div>
                <h3 className="item-name">{item.name}</h3>
                <p className="item-price">
                  ${item.price.toFixed(2)}
                </p>
              </div>
            ))}
          </div>
        )}
        
        {!loading && items.length === 0 && (
          <div className="empty-state">
            <p className="empty-text">No items found.</p>
            <p className="empty-subtext">Create one to get started.</p>
          </div>
        )}
      </div>

      {isModalOpen && (
        <div className="modal-overlay">
          <div className="glass-panel modal-content">
            <div className="modal-header">
              <h2 className="modal-title">
                {currentItem ? 'Edit Item' : 'New Item'}
              </h2>
              <button 
                onClick={() => setIsModalOpen(false)}
                className="close-btn"
              >
                <X size={24} />
              </button>
            </div>
            
            <form onSubmit={handleSave}>
              <div className="form-group">
                <div>
                  <label className="label">Item Name</label>
                  <input 
                    name="name" 
                    defaultValue={currentItem?.name}
                    className="input-field"
                    placeholder="e.g. Premium Widget"
                    required
                  />
                </div>
                <div>
                  <label className="label">Price ($)</label>
                  <input 
                    name="price" 
                    type="number" 
                    step="0.01"
                    defaultValue={currentItem?.price}
                    className="input-field"
                    placeholder="0.00"
                    required
                  />
                </div>
              </div>
              
              <div className="form-actions">
                <button 
                  type="button"
                  onClick={() => setIsModalOpen(false)}
                  className="btn"
                >
                  Cancel
                </button>
                <button 
                  type="submit"
                  className="btn btn-primary"
                >
                  <Save size={18} />
                  {currentItem ? 'Save Changes' : 'Create Item'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}

export default App
