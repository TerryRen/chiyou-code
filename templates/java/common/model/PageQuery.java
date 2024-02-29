package {xxx}.common.model;

import io.swagger.v3.oas.annotations.media.Schema;

import java.io.Serializable;

/**
 * 分页查询参数
 *
 * @author mmc
 */
public class PageQuery implements Serializable {
    private static final long serialVersionUID = -231509528616988497L;
    /**
     * 分页-页数
     **/
    @Schema(description = "分页-页数, 默认-1", example = "1")
    private Integer page = 1;

    /**
     * 分页-每页大小
     **/
    @Schema(description = "分页-每页大小, 默认-10", example = "10")
    private Integer size = 10;

    public Integer getPage() {
        return page;
    }

    public void setPage(Integer page) {
        this.page = page;
    }

    public Integer getSize() {
        return size;
    }

    public void setSize(Integer size) {
        this.size = size;
    }
}

