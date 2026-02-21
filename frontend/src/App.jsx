import { useState, useEffect, useCallback } from 'react'
import { Plus, Trash2, Pencil, X, Save, RefreshCw, ChevronLeft, ChevronRight, AlertCircle } from 'lucide-react'

const PAGE_SIZE = 20

function App() {
  const [items, setItems] = useState([])
  const [loading, setLoading] = useState(true)
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [currentItem, setCurrentItem] = useState(null)
  const [total, setTotal] = useState(0)
  const [offset, setOffset] = useState(0)
  const [toast, setToast] = useState(null)

  const API_URL = import.meta.env.VITE_API_URL || '/api'

  const showToast = (message, type = 'error') => {
    setToast({ message, type })
    setTimeout(() => setToast(null), 4000)
  }

  const fetchItems = useCallback(async () => {
    setLoading(true)
    try {
      const res = await fetch(`${API_URL}/items?limit=${PAGE_SIZE}&offset=${offset}`)
      if (!res.ok) throw new Error('Failed to fetch items')
      const data = await res.json()
      setItems(data.items || [])
      setTotal(data.total || 0)
    } catch (error) {
      showToast(error.message)
    } finally {
      setLoading(false)
    }
  }, [API_URL, offset])

  useEffect(() => {
    fetchItems()
  }, [fetchItems])

  const handleDelete = async (id) => {
    if (!confirm('Are you sure you want to delete this item?')) return
    try {
      const res = await fetch(`${API_URL}/items/${id}`, { method: 'DELETE' })
      if (!res.ok) throw new Error('Failed to delete item')
      showToast('Item deleted', 'success')
      fetchItems()
    } catch (error) {
      showToast(error.message)
    }
  }

  const handleSave = async (e) => {
    e.preventDefault()
    const formData = new FormData(e.target)
    const name = formData.get('name').trim()
    const price = parseFloat(formData.get('price'))

    if (!name) {
      showToast('Name is required')
      return
    }
    if (name.length > 255) {
      showToast('Name must be 255 characters or less')
      return
    }
    if (isNaN(price) || price < 0) {
      showToast('Price must be a non-negative number')
      return
    }

    const payload = { name, price }

    try {
      let res
      if (currentItem) {
        res = await fetch(`${API_URL}/items/${currentItem.id}`, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload)
        })
      } else {
        res = await fetch(`${API_URL}/items`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload)
        })
      }
      if (!res.ok) {
        const text = await res.text()
        throw new Error(text || 'Failed to save item')
      }
      showToast(currentItem ? 'Item updated' : 'Item created', 'success')
      setIsModalOpen(false)
      setCurrentItem(null)
      fetchItems()
    } catch (error) {
      showToast(error.message)
    }
  }

  const openModal = (item = null) => {
    setCurrentItem(item)
    setIsModalOpen(true)
  }

  const totalPages = Math.ceil(total / PAGE_SIZE)
  const currentPage = Math.floor(offset / PAGE_SIZE) + 1

  return (
    <div className="app-layout">
      {toast && (
        <div className={`toast toast-${toast.type}`}>
          <AlertCircle size={16} />
          {toast.message}
        </div>
      )}

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
          <>
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

            {totalPages > 1 && (
              <div className="pagination">
                <button
                  className="btn"
                  disabled={currentPage <= 1}
                  onClick={() => setOffset(offset - PAGE_SIZE)}
                >
                  <ChevronLeft size={18} />
                  Previous
                </button>
                <span className="pagination-info">
                  Page {currentPage} of {totalPages} ({total} items)
                </span>
                <button
                  className="btn"
                  disabled={currentPage >= totalPages}
                  onClick={() => setOffset(offset + PAGE_SIZE)}
                >
                  Next
                  <ChevronRight size={18} />
                </button>
              </div>
            )}
          </>
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
                    maxLength={255}
                    required
                  />
                </div>
                <div>
                  <label className="label">Price ($)</label>
                  <input
                    name="price"
                    type="number"
                    step="0.01"
                    min="0"
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
