package {xxx}.common.model;

import com.newegg.logistics.toa.common.context.CurrentUserContext;
import io.swagger.v3.oas.annotations.media.Schema;
import lombok.Data;
import lombok.experimental.Accessors;

import java.io.Serializable;

/**
 * @author mmc
 */
@Data
@Accessors(chain = true)
public class BaseModel implements Serializable {
    @Schema(description = "Key", example = "1")
    private Integer transactionNumber;

    /**
     * 数据修改版本号 (用于乐观并发控制)
     **/
    @Schema(description = "数据修改版本号", example = "0")
    private Integer version;

    @Schema(description = "创建人", example = "admin")
    private String inUser;

    @Schema(description = "创建时间", example = "1622103236423")
    private Long inDate;

    @Schema(description = "最后修改人", example = "admin")
    private String lastEditUser;

    @Schema(description = "最后修改时间", example = "1622103236423")
    private Long lastEditDate;

    @Schema(description = "是否删除", example = "false")
    private Boolean deleted;

    /**
     * 内部单表排序使用
     */
    private String sortColumns;

    /**
     * 填充insert时的通用字段
     *
     * @return 本对象，以供链式调用
     */
    public BaseModel fillInsertColumn() {
        long now = System.currentTimeMillis();
        this.setVersion(0);
        this.setDeleted(Boolean.FALSE);
        this.setInDate(now);
        this.setInUser(CurrentUserContext.getCurrentUser().getUserEmail());
        this.setLastEditDate(now);
        this.setLastEditUser(CurrentUserContext.getCurrentUser().getUserEmail());
        return this;
    }

    /**
     * 填充update时的通用字段
     *
     * @return 本对象，以供链式调用
     */
    public BaseModel fillUpdateColumn() {
        this.setLastEditDate(System.currentTimeMillis());
        this.setLastEditUser(CurrentUserContext.getCurrentUser().getUserEmail());
        return this;
    }
}
