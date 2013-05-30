package petadb

type PagedInfo struct {
	CurrentPageIndex int
	PageSize         int
	TotalItemCount   int
	TotalPageCount   int
}
