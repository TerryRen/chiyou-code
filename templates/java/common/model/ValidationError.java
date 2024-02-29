package {xxx}.common.model;

import io.swagger.v3.oas.annotations.media.Schema;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * @author mmc
 */
@Data
@NoArgsConstructor
@AllArgsConstructor
@Schema(description = "字段校验数据")
public class ValidationError {
    @Schema(description = "字段名")
    private String propertyName;
    
    @Schema(description = "错误信息")
    private String errorMessage;
}
