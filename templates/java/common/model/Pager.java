package {xxx}.common.model;

import io.swagger.v3.oas.annotations.media.Schema;

import java.io.Serializable;
import java.util.ArrayList;
import java.util.List;

/**
 * @author mmc
 */
@Schema(description = "分页对象")
public class Pager<T> implements Serializable {
    /**
     * serialVersionUID
     */
    private static final long serialVersionUID = -1886255970252256895L;

    /**
     * 当前页号
     */
    @Schema(description = "当前页号", example = "1")
    private int currentPage;

    @Schema(description = "开始数量", example = "1")
    private int startNum;

    @Schema(description = "结束数量", example = "1")
    private int endNum;

    /**
     * 每页记录数
     */
    @Schema(description = "每页记录数", example = "10")
    private int pageSize;

    /**
     * 结果集起始记录数
     */
    @Schema(description = "起始记录数", example = "0")
    private int startRow;

    /**
     * 总计页数
     */
    @Schema(description = "总计页数", example = "1")
    private int totalPages;

    /**
     * 总计记录数
     */
    @Schema(description = "总计记录数", example = "100")
    private int total;

    /**
     * 结果集
     */
    private List<T> results = new ArrayList<>();

    public Pager() {
        currentPage = 1;
        startRow = 0;
        total = 0;
        startNum = 0;
        endNum = 10;
    }

    /**
     * 构造函数 通过记录总数和每页记录数计算出相关信息
     *
     * @param total      int //记录总数
     * @param onePageSize int //每页记录数
     */
    public Pager(int total, int onePageSize) {
        this.pageSize = onePageSize;
        this.total = total;
        totalPages = this.total / pageSize;
        int mod = this.total % pageSize;
        if (mod > 0) {
            totalPages++;
        }
        if (this.total == 0) {
            currentPage = 0;
        } else {
            currentPage = 1;
        }
        startRow = 0;
        if (this.total == 0) {
            startNum = 0;
        } else {
            startNum = 1;
        }
        if (this.total <= pageSize) {
            endNum = this.total;
        } else {
            if (currentPage == totalPages) {
                endNum = this.total;
            } else {
                endNum = pageSize * currentPage;
            }
        }
    }

    /**
     * 构造函数 通过记录总数和每页记录数计算出相关信息
     *
     * @param total      int //记录总数
     * @param onePageSize int //每页记录数
     */
    public Pager(int total, int onePageSize, int currentPage) {
        this.pageSize = onePageSize;
        this.total = total;
        totalPages = this.total / pageSize;

        int mod = this.total % pageSize;
        if (mod > 0) {
            totalPages++;
        }
        if (totalPages < currentPage) {
            currentPage = totalPages;
        }
        if (currentPage < 1) {
            currentPage = 1;
        }
        if (this.total == 0) {
            this.currentPage = 1;
        } else {
            this.currentPage = currentPage;
        }
        startRow = 0;
        if (this.total == 0) {
            startNum = 0;
        } else {
            startNum = (currentPage - 1) * onePageSize;
        }
        if (this.total <= pageSize) {
            endNum = this.total;
        } else {
            if (currentPage == totalPages) {
                endNum = this.total;
            } else {
                endNum = pageSize * currentPage;
            }
        }
    }

    /**
     * 第一页
     */
    public Pager<T> first() {
        if (total == 0) {
            currentPage = 0;
        } else {
            currentPage = 1;
        }
        startRow = 0;
        if (total == 0) {
            startNum = 0;
        } else {
            startNum = 1;
        }
        if (totalPages <= pageSize) {
            endNum = total;
        } else {
            if (currentPage == totalPages) {
                endNum = total;
            } else {
                endNum = pageSize * currentPage;
            }
        }
        return this;
    }

    /**
     * 最后一页
     */
    public Pager<T> last() {
        if (totalPages != 0) {
            currentPage = totalPages;
            startRow = (currentPage - 1) * pageSize;
            startNum = startRow + 1;
            endNum = total;
        }
        return this;
    }

    /**
     * 下一页
     */
    public Pager<T> next() {

        if (currentPage < totalPages && totalPages != 0) {
            currentPage++;
        }
        if (totalPages != 0) {
            startRow = (currentPage - 1) * pageSize;
        }
        startNum = startRow + 1;
        if (currentPage < totalPages && totalPages != 0) {
            endNum = startRow + pageSize;
        } else {
            endNum = total;
        }
        return this;
    }

    /**
     * 上一页
     */
    public Pager<T> previous() {
        if (currentPage <= 1) {
            first();
        } else {
            currentPage--;
            startRow = (currentPage - 1) * pageSize;
            startNum = startRow + 1;
            if (currentPage < totalPages) {
                endNum = startRow + pageSize;
            } else {
                endNum = total;
            }
        }
        return this;
    }

    /**
     * 更新页
     *
     * @param currentPage int
     */
    public Pager<T> refresh(int currentPage) {
        this.currentPage = currentPage;
        if (this.currentPage > totalPages) {
            last();
        }
        if (this.currentPage < 0) {
            first();
        }
        return this;
    }

    /**
     * 指定页码
     *
     * @param pageNo int
     */
    public Pager<T> go(int pageNo) {
        if (pageNo <= 1) {
            first();
        } else if (pageNo >= totalPages) {
            last();
            if (totalPages < 1) {
                currentPage = 0;
            }
        } else {
            startRow = (currentPage - 1) * pageSize;
            startNum = startRow + 1;
            endNum = startRow + pageSize;
        }
        return this;
    }

    public int getCurrentPage() {
        return currentPage;
    }

    public int getPageSize() {
        return pageSize;
    }

    public int getStartRow() {
        return startRow;
    }

    public int getTotalPages() {
        return totalPages;
    }

    public int getTotal() {
        return total;
    }

    public Pager<T> setCurrentPage(int currentPage) {
        this.currentPage = currentPage;
        return this;
    }

    public Pager<T> setPageSize(int pageSize) {
        this.pageSize = pageSize;
        return this;
    }

    public Pager<T> setStartRow(int startRow) {
        this.startRow = startRow;
        return this;
    }

    public Pager<T> setTotalPages(int totalPages) {
        this.totalPages = totalPages;
        return this;
    }

    public Pager<T> setTotal(int total) {
        this.total = total;
        return this;
    }

    public int getEndNum() {
        return endNum;
    }

    public Pager<T> setEndNum(int endNum) {
        this.endNum = endNum;
        return this;
    }

    public int getStartNum() {
        if (startNum < 0) {
            startNum = 0;
        }
        return startNum;
    }

    public Pager<T> setStartNum(int startNum) {
        this.startNum = startNum;
        return this;
    }

    public List<T> getResults() {
        return results;
    }

    public Pager<T> setResults(List<T> results) {
        this.results = results;
        return this;
    }
}

