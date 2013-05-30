package petadb

type PagedList struct {
	CurrentPageIndex int
	PageSize         int
	List             interface{}
	TotalItemCount   int
	TotalPageCount   int
}
